type PlaceholderPageProps = {
  eyebrow: string;
  title: string;
  description: string;
};

export function PlaceholderPage({
  eyebrow,
  title,
  description,
}: PlaceholderPageProps) {
  return (
    <section className="page">
      <p className="eyebrow">{eyebrow}</p>
      <h1>{title}</h1>
      <p className="page-description">{description}</p>
      <div className="placeholder-panel">
        <span className="placeholder-icon" aria-hidden="true">
          +
        </span>
        <p>Feature content will be added here.</p>
      </div>
    </section>
  );
}
