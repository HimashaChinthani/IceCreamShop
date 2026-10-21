import {
  LayoutDashboard,
  Package,
  Users,
  LogOut,
  IceCream,
} from "lucide-react";
import { NavLink, useNavigate } from "react-router-dom";
import { logout } from "../services/auth";

const menuItems = [
  { label: "Dashboard", icon: LayoutDashboard, path: "/" },
  { label: "Inventory", icon: Package, path: "/inventory" },
  { label: "Customers", icon: Users, path: "/customers" },
];

export default function Sidebar() {
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login", { replace: true });
  }

  return (
    <aside className="sidebar">
      <div className="brand">
        <IceCream size={30} />
        <span>IceCream POS</span>
      </div>

      <p className="menu-heading">WORKSPACE</p>

      <nav className="menu">
        {menuItems.map(({ label, icon: Icon, path }) => (
          <NavLink
            className={({ isActive }) => `menu-item ${isActive ? "active" : ""}`}
            key={label}
            to={path}
          >
            <Icon size={19} />
            <span>{label}</span>
          </NavLink>
        ))}
      </nav>

      <div className="sidebar-bottom">
        <button className="menu-item" onClick={handleLogout}>
          <LogOut size={19} />
          <span>Logout</span>
        </button>
      </div>
    </aside>
  );
}