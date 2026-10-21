import { useEffect, useState } from "react";
import {
  Package,
  Boxes,
  AlertTriangle,
  DollarSign,
} from "lucide-react";
import api from "../services/api";
import { getCurrentUser } from "../services/auth";

const statDefinitions = [
  {
    title: "Total Products",
    value: "0",
    description: "Products in inventory",
    icon: Package,
  },
  {
    title: "Total Stock",
    value: "0",
    description: "Units available",
    icon: Boxes,
  },
  {
    title: "Low Stock Items",
    value: "0",
    description: "Items needing attention",
    icon: AlertTriangle,
  },
  {
    title: "Inventory Value",
    value: "Rs. 0.00",
    description: "Estimated stock value",
    icon: DollarSign,
  },
];

export default function Dashboard() {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const user = getCurrentUser();

  useEffect(() => {
    api.get("/inventory/")
      .then((response) => setItems(response.data))
      .catch((err) => setError(err.response?.data?.error || "Unable to load dashboard data."))
      .finally(() => setLoading(false));
  }, []);

  const stats = [
    { ...statDefinitions[0], value: items.length },
    { ...statDefinitions[1], value: items.reduce((sum, item) => sum + Number(item.quantity || 0), 0) },
    { ...statDefinitions[2], value: items.filter((item) => Number(item.quantity) <= 10).length },
    { ...statDefinitions[3], value: `Rs. ${items.reduce((sum, item) => sum + Number(item.quantity || 0) * Number(item.price || 0), 0).toFixed(2)}` },
  ];

  return (
    <main className="dashboard">
      <header className="page-header">
        <div>
          <p className="eyebrow">OVERVIEW</p>
          <h1>Dashboard</h1>
          <p className="subtitle">
            Welcome back! Here's your inventory overview.
          </p>
        </div>

        <div className="user-badge">
          <div className="avatar">{user?.full_name?.charAt(0) || "U"}</div>
          <div>
            <strong>{user?.full_name || "User"}</strong>
            <p>{user?.role || "Staff"}</p>
          </div>
        </div>
      </header>

      <section className="stats-grid">
        {stats.map(({ title, value, description, icon: Icon }) => (
          <article className="stat-card" key={title}>
            <div className="stat-top">
              <span>{title}</span>
              <div className="stat-icon">
                <Icon size={21} />
              </div>
            </div>

            <h2>{value}</h2>
            <p>{description}</p>
          </article>
        ))}
      </section>

      {error && <div className="notice error-notice">{error}</div>}
      <section className="content-card">
        <div>
          <h2>Inventory Overview</h2>
          <p>{loading ? "Loading inventory data..." : "Your latest inventory summary."}</p>
        </div>

        <div className="empty-state">
          <Package size={42} />
          <h3>{items.length ? `${items.length} products currently tracked` : "No inventory displayed yet"}</h3>
          <p>Use the Inventory page to manage your products and stock.</p>
        </div>
      </section>
    </main>
  );
}