
import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/AppLayout";

// Applications Management — Issue #7
// Applications Management — Issue #7
import Applications from "./APPLICATION/Applications";
import Matcher from "./MATCHER/Matcher";
// Dashboard — Issue #9
import Dashboard from "./APPLICATION/Dashboard";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<Navigate replace to="/dashboard" />} />
        <Route
          path="/dashboard"
          element={
            <Dashboard />
          }
        />

        {/* Applications Management — Issue #7 */}
        <Route
          path="/applications"
          element={<Applications />}
        />

        <Route
          path="/keyword-matcher"
          element={<Matcher />}
        />
        <Route
          path="*"
          element={<Navigate replace to="/dashboard" />}
        />
      </Route>
    </Routes>
  );
}
