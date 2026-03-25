import { NavLink, Outlet } from "react-router-dom";

const links = [
  { to: "/", label: "Rounds" },
  { to: "/promotions", label: "Promotions" }
];

export function Layout() {
  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand">
          <p className="eyebrow">Clawbot Platform</p>
          <h1>Trust Lab Operator</h1>
          <p className="muted">Review rounds, promotions, reports, and detector behavior.</p>
        </div>
        <nav className="nav">
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              end={link.to === "/"}
              className={({ isActive }) => (isActive ? "nav-link nav-link-active" : "nav-link")}
            >
              {link.label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
