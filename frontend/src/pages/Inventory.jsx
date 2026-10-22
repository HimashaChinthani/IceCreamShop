import { useEffect, useState } from "react";
import {
  Search,
  Plus,
  Pencil,
  Trash2,
  Package,
  X,
} from "lucide-react";
import api from "../services/api";
import { getCurrentUser } from "../services/auth";

const emptyForm = {
  name: "",
  category: "",
  sku: "",
  quantity: 0,
  price: 0,
};

export default function Inventory() {
  const [items, setItems] = useState([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [modalOpen, setModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState(null);
  const [form, setForm] = useState(emptyForm);
  const [saving, setSaving] = useState(false);

  const user = getCurrentUser();
  const role = user?.role?.toLowerCase();

  const canCreate = role === "owner" || role === "manager";
  const canEdit = ["owner", "manager", "cashier"].includes(role);
  const canDelete = role === "owner";

  async function loadInventory() {
    setError("");

    try {
      const response = await api.get("/inventory/");
      setItems(response.data);
    } catch (err) {
      setError(
        err.response?.data?.error ||
          "Unable to load inventory. Check your API connection."
      );
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadInventory();
  }, []);

  function openCreateModal() {
    setEditingItem(null);
    setForm(emptyForm);
    setError("");
    setSuccess("");
    setModalOpen(true);
  }

  function openEditModal(item) {
    setEditingItem(item);
    setForm({
      name: item.name,
      category: item.category || "",
      sku: item.sku,
      quantity: item.quantity,
      price: item.price,
    });
    setError("");
    setSuccess("");
    setModalOpen(true);
  }

  function updateField(event) {
    const { name, value } = event.target;

    setForm((previous) => ({
      ...previous,
      [name]: value,
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    setSuccess("");

    const payload = {
      name: form.name.trim(),
      category: form.category.trim(),
      sku: form.sku.trim(),
      quantity: Number(form.quantity),
      price: Number(form.price),
    };

    try {
      if (editingItem) {
        await api.put(`/inventory/${editingItem.id}`, payload);
        setSuccess("Inventory item updated successfully.");
      } else {
        await api.post("/inventory/", payload);
        setSuccess("Inventory item created successfully.");
      }

      setModalOpen(false);
      await loadInventory();
    } catch (err) {
      setError(
        err.response?.data?.error ||
          "Unable to save the inventory item."
      );
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(item) {
    const confirmed = window.confirm(
      `Delete "${item.name}" from inventory?`
    );

    if (!confirmed) return;

    setError("");
    setSuccess("");

    try {
      await api.delete(`/inventory/${item.id}`);
      setItems((previous) =>
        previous.filter((entry) => entry.id !== item.id)
      );
      setSuccess("Inventory item deleted successfully.");
    } catch (err) {
      setError(
        err.response?.data?.error ||
          "Unable to delete the inventory item."
      );
    }
  }

  const filteredItems = items.filter((item) => {
    const query = search.toLowerCase();

    return (
      item.name.toLowerCase().includes(query) ||
      item.category?.toLowerCase().includes(query) ||
      item.sku.toLowerCase().includes(query)
    );
  });

  return (
    <main className="dashboard">
      <header className="page-header">
        <div>
          <p className="eyebrow">PRODUCT MANAGEMENT</p>
          <h1>Inventory</h1>
          <p className="subtitle">
            Manage your products and stock levels.
          </p>
        </div>

        {canCreate && (
          <button className="primary-button" onClick={openCreateModal}>
            <Plus size={18} />
            Add Product
          </button>
        )}
      </header>

      {error && <div className="notice error-notice">{error}</div>}
      {success && (
        <div className="notice success-notice">{success}</div>
      )}

      <section className="content-card inventory-card">
        <div className="inventory-toolbar">
          <div>
            <h2>All Products</h2>
            <p>{items.length} products in inventory</p>
          </div>

          <div className="search-box">
            <Search size={18} />
            <input
              type="search"
              placeholder="Search products or SKU..."
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </div>
        </div>

        {loading ? (
          <p className="inventory-message">Loading inventory...</p>
        ) : (
          <div className="table-wrapper">
            <table className="inventory-table">
              <thead>
                <tr>
                  <th>Product</th>
                  <th>SKU</th>
                  <th>Category</th>
                  <th>Quantity</th>
                  <th>Price</th>
                  <th>Stock Value</th>
                  {(canEdit || canDelete) && <th>Actions</th>}
                </tr>
              </thead>

              <tbody>
                {filteredItems.map((item) => (
                  <tr key={item.id}>
                    <td>
                      <div className="product-cell">
                        <div className="product-icon">
                          <Package size={20} />
                        </div>
                        <strong>{item.name}</strong>
                      </div>
                    </td>
                    <td>{item.sku}</td>
                    <td>{item.category}</td>
                    <td>
                      <span
                        className={
                          item.quantity <= 10
                            ? "stock-badge low-stock"
                            : "stock-badge"
                        }
                      >
                        {item.quantity}
                      </span>
                    </td>
                    <td>
                      Rs. {Number(item.price).toFixed(2)}
                    </td>
                    <td>
                      Rs.{" "}
                      {(
                        Number(item.quantity) * Number(item.price)
                      ).toFixed(2)}
                    </td>
                    {(canEdit || canDelete) && (
                      <td>
                        <div className="action-buttons">
                          {canEdit && (
                            <button
                              className="icon-button"
                              title="Edit product"
                              onClick={() => openEditModal(item)}
                            >
                              <Pencil size={17} />
                            </button>
                          )}

                          {canDelete && (
                            <button
                              className="icon-button delete-button"
                              title="Delete product"
                              onClick={() => handleDelete(item)}
                            >
                              <Trash2 size={17} />
                            </button>
                          )}
                        </div>
                      </td>
                    )}
                  </tr>
                ))}

                {filteredItems.length === 0 && (
                  <tr>
                    <td
                      colSpan={canEdit || canDelete ? 7 : 6}
                      className="empty-table"
                    >
                      No matching inventory items found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {modalOpen && (
        <div className="modal-backdrop">
          <section className="inventory-modal">
            <header className="modal-header">
              <h2>
                {editingItem ? "Edit Product" : "Add Product"}
              </h2>

              <button
                className="icon-button"
                onClick={() => setModalOpen(false)}
                aria-label="Close modal"
              >
                <X size={20} />
              </button>
            </header>

            <form className="inventory-form" onSubmit={handleSubmit}>
              <label>
                Product name
                <input
                  name="name"
                  value={form.name}
                  onChange={updateField}
                  required
                  maxLength={120}
                />
              </label>
              <label>
                Category
                <input
                  name="category"
                  value={form.category}
                  onChange={updateField}
                  placeholder="Chocolate"
                  required
                />
              </label>

              <label>
                SKU
                <input
                  name="sku"
                  value={form.sku}
                  onChange={updateField}
                  required
                  maxLength={50}
                />
              </label>

              <label>
                Quantity
                <input
                  name="quantity"
                  type="number"
                  min="0"
                  step="1"
                  value={form.quantity}
                  onChange={updateField}
                  required
                />
              </label>

              <label>
                Price (LKR)
                <input
                  name="price"
                  type="number"
                  min="0"
                  step="0.01"
                  value={form.price}
                  onChange={updateField}
                  required
                />
              </label>

              {error && (
                <p className="login-error">{error}</p>
              )}

              <button
                className="primary-button modal-submit"
                type="submit"
                disabled={saving}
              >
                {saving
                  ? "Saving..."
                  : editingItem
                    ? "Save Changes"
                    : "Create Product"}
              </button>
            </form>
          </section>
        </div>
      )}
    </main>
  );
}