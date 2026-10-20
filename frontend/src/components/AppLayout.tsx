import { NavLink, Outlet } from "react-router-dom";

const navigation = [
  { label: "Dashboard", to: "/dashboard" },
  { label: "Applications", to: "/applications" },
  { label: "Keyword Matcher", to: "/keyword-matcher" },
];

export function AppLayout() {
  return (
    <div className="app-shell">
      <header className="topbar">
        <a className="brand" href="/dashboard" aria-label="JobTrack home">
          <span className="brand-mark" aria-hidden="true">
            J
          </span>
          <span>JobTrack</span>
        </a>
        <nav className="primary-nav" aria-label="Main navigation">
          {navigation.map((item) => (
            <NavLink
              key={item.to}
              className={({ isActive }) =>
                isActive ? "nav-link nav-link-active" : "nav-link"
              }
              to={item.to}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
        <span className="topbar-note">Your job search, in one place</span>
      </header>
      <main className="main-content">
        <Outlet />
      </main>
      <footer className="footer">
        <span>JobTrack</span>
        <span>Make your next move with confidence.</span>
      </footer>
    </div>
  );
}
