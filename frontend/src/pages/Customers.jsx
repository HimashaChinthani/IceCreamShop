import { useCallback, useEffect, useState } from "react";
import { Pencil, Plus, Search, Trash2, UserRound, X } from "lucide-react";
import api from "../services/api";
import { getCurrentUser } from "../services/auth";

const emptyForm = { name: "", email: "", phone: "", address: "" };

export default function Customers() {
  const [customers, setCustomers] = useState([]);
  const [search, setSearch] = useState("");
  const [form, setForm] = useState(emptyForm);
  const [editingCustomer, setEditingCustomer] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const role = getCurrentUser()?.role?.toLowerCase();
  const canEdit = ["owner", "manager", "cashier"].includes(role);
  const canDelete = role === "owner";

  const loadCustomers = useCallback(async () => {
    setLoading(true);
    try {
      const response = await api.get("/customers/", {
        params: search ? { search } : undefined,
      });
      setCustomers(response.data);
      setError("");
    } catch (err) {
      setError(err.response?.data?.error || "Unable to load customers.");
    } finally {
      setLoading(false);
    }
  }, [search]);

  useEffect(() => {
    const timer = setTimeout(() => {
      loadCustomers();
    }, 250);
    return () => clearTimeout(timer);
  }, [loadCustomers]);

  function openModal(customer = null) {
    setEditingCustomer(customer);
    setForm(customer ? {
      name: customer.name,
      email: customer.email || "",
      phone: customer.phone || "",
      address: customer.address || "",
    } : emptyForm);
    setError("");
    setSuccess("");
    setModalOpen(true);
  }

  function updateField(event) {
    setForm((previous) => ({ ...previous, [event.target.name]: event.target.value }));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setSaving(true);
    setError("");
    try {
      if (editingCustomer) {
        await api.put(`/customers/${editingCustomer.id}`, form);
        setSuccess("Customer updated successfully.");
      } else {
        await api.post("/customers/", form);
        setSuccess("Customer created successfully.");
      }
      setModalOpen(false);
      await loadCustomers();
    } catch (err) {
      setError(err.response?.data?.error || "Unable to save customer.");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(customer) {
    if (!window.confirm(`Delete "${customer.name}"?`)) return;
    try {
      await api.delete(`/customers/${customer.id}`);
      setCustomers((current) => current.filter((item) => item.id !== customer.id));
      setSuccess("Customer deleted successfully.");
    } catch (err) {
      setError(err.response?.data?.error || "Unable to delete customer.");
    }
  }

  return (
    <main className="dashboard">
      <header className="page-header">
        <div>
          <p className="eyebrow">CUSTOMER MANAGEMENT</p>
          <h1>Customers</h1>
          <p className="subtitle">Manage your customer records.</p>
        </div>
        {canEdit && <button className="primary-button" onClick={() => openModal()}>
          <Plus size={18} /> Add Customer
        </button>}
      </header>
      {error && <div className="notice error-notice">{error}</div>}
      {success && <div className="notice success-notice">{success}</div>}
      <section className="content-card">
        <div className="inventory-toolbar">
          <div><h2>All Customers</h2><p>{customers.length} customer records</p></div>
          <div className="search-box"><Search size={18} /><input type="search" placeholder="Search customers..." value={search} onChange={(event) => setSearch(event.target.value)} /></div>
        </div>
        {loading ? <p className="inventory-message">Loading customers...</p> : (
          <div className="table-wrapper">
            <table className="inventory-table">
              <thead><tr><th>Customer</th><th>Email</th><th>Phone</th><th>Address</th>{canEdit && <th>Actions</th>}</tr></thead>
              <tbody>
                {customers.map((customer) => <tr key={customer.id}>
                  <td><div className="product-cell"><div className="product-icon"><UserRound size={20} /></div><strong>{customer.name}</strong></div></td>
                  <td>{customer.email || "—"}</td><td>{customer.phone || "—"}</td><td>{customer.address || "—"}</td>
                  {canEdit && <td><div className="action-buttons">
                    <button className="icon-button" title="Edit customer" onClick={() => openModal(customer)}><Pencil size={17} /></button>
                    {canDelete && <button className="icon-button delete-button" title="Delete customer" onClick={() => handleDelete(customer)}><Trash2 size={17} /></button>}
                  </div></td>}
                </tr>)}
              </tbody>
            </table>
            {!customers.length && <p className="empty-table">No customers found.</p>}
          </div>
        )}
      </section>
      {modalOpen && <div className="modal-backdrop">
        <section className="inventory-modal">
          <div className="modal-header"><h2>{editingCustomer ? "Edit Customer" : "Add Customer"}</h2><button className="icon-button" onClick={() => setModalOpen(false)}><X size={18} /></button></div>
          <form className="inventory-form" onSubmit={handleSubmit}>
            <label>Name<input name="name" value={form.name} onChange={updateField} required /></label>
            <label>Email<input type="email" name="email" value={form.email} onChange={updateField} /></label>
            <label>Phone<input name="phone" value={form.phone} onChange={updateField} /></label>
            <label>Address<input name="address" value={form.address} onChange={updateField} /></label>
            <button className="primary-button modal-submit" disabled={saving}>{saving ? "Saving..." : "Save Customer"}</button>
          </form>
        </section>
      </div>}
    </main>
  );
}
