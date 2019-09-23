import React, { useEffect } from "react";

const Index: React.FC = () => {
    useEffect(() => {
        (async () => {
            try {
                const res = await fetch("/api/index");
                if (!res.ok) {
                    throw new Error(res.statusText);
                }
                const data = await res.json();
                window.console.log(data);
            } catch (err) {
                window.console.error(err.message);
            }
        })();
    });
    return (
      <div>index</div>
    );
};

export default Index;
