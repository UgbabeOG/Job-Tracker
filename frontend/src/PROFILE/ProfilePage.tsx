import { useEffect, useState } from "react";
import { getProfile, updateProfile } from "./profile";
import { Button, LoadingState, ErrorMessage } from "../components/ui";

export default function ProfilePage() {
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [linkedin, setLinkedin] = useState("");
  const [portfolio, setPortfolio] = useState("");
  const [cvText, setCvText] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    getProfile()
      .then((p) => {
        if (p) {
          setFullName(p.full_name || "");
          setEmail(p.email || "");
          setPhone(p.phone || "");
          setLinkedin(p.linkedin_url || "");
          setPortfolio(p.portfolio_url || "");
          setCvText(p.cv_text || "");
        }
      })
      .catch(() => setError("Failed to load profile"))
      .finally(() => setLoading(false));
  }, []);

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSaved(false);
    try {
      await updateProfile({
        full_name: fullName || null,
        email: email || null,
        phone: phone || null,
        linkedin_url: linkedin || null,
        portfolio_url: portfolio || null,
        cv_text: cvText || null,
      });
      setSaved(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save");
    } finally {
      setSaving(false);
    }
  }

  if (loading) return <LoadingState />;

  return (
    <section className="applications-page">
      <header className="applications-header">
        <div>
          <h1>Profile</h1>
          <p>Store your reusable details and CV for the Apply Assistant.</p>
        </div>
      </header>
      {error && <ErrorMessage message={error} />}
      {saved && <p>Saved.</p>}
      <form onSubmit={handleSave}>
        <div style={{ display: "grid", gap: "1rem", maxWidth: "600px" }}>
          <label>
            Full Name
            <input value={fullName} onChange={(e) => setFullName(e.target.value)} />
          </label>
          <label>
            Email
            <input value={email} onChange={(e) => setEmail(e.target.value)} />
          </label>
          <label>
            Phone
            <input value={phone} onChange={(e) => setPhone(e.target.value)} />
          </label>
          <label>
            LinkedIn URL
            <input value={linkedin} onChange={(e) => setLinkedin(e.target.value)} />
          </label>
          <label>
            Portfolio URL
            <input value={portfolio} onChange={(e) => setPortfolio(e.target.value)} />
          </label>
          <label>
            CV Text
            <textarea value={cvText} onChange={(e) => setCvText(e.target.value)} rows={10} />
          </label>
          <Button type="submit" disabled={saving}>{saving ? "Saving..." : "Save"}</Button>
        </div>
      </form>
    </section>
  );
}
