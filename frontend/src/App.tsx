
import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/AppLayout";
import { PlaceholderPage } from "./components/PlaceholderPage";

// Applications Management — Issue #7
// Applications Management — Issue #7
import Applications from "./APPLICATION/Applications";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<Navigate replace to="/dashboard" />} />
        <Route
          path="/dashboard"
          element={
            <PlaceholderPage
              eyebrow="Overview"
              title="Dashboard"
              description="Your job-search overview will appear here once dashboard features are connected."
            />
          }
        />

        {/* Applications Management — Issue #7 */}
        <Route
          path="/applications"
          element={<Applications />}
        />

        <Route
          path="/keyword-matcher"
          element={
            <PlaceholderPage
              eyebrow="Prepare your application"
              title="Keyword Matcher"
              description="Compare CV text with a job description here when keyword matching is implemented."
            />
          }
        />
        <Route
          path="*"
          element={<Navigate replace to="/dashboard" />}
        />
      </Route>
    </Routes>
  );
}
