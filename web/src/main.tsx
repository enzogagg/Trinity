import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { Health, isHealthy } from "./health";

function App() {
  const [ok, setOk] = useState<boolean | null>(null);
  useEffect(() => {
    fetch("/api/health")
      .then((r) => r.json() as Promise<Health>)
      .then((h) => setOk(isHealthy(h)))
      .catch(() => setOk(false));
  }, []);
  return <h1>Trinity — API: {ok === null ? "..." : ok ? "OK" : "KO"}</h1>;
}

createRoot(document.getElementById("root")!).render(<App />);
