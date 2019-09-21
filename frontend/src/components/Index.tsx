import React, { useCallback, useState } from "react";
import firebase from "firebase/app";
import "firebase/auth";

const Index: React.FC = () => {
    const [submitting, setSubmitting] = useState(false);
    const onClickSignIn = useCallback(() => {
        setSubmitting(true);
        (async () => {
            try {
                const provider = new firebase.auth.GoogleAuthProvider();
                const credential: firebase.auth.UserCredential = await firebase.auth().signInWithPopup(provider);
                if (!credential.user) {
                    throw new Error("no credential user");
                }
                const token = await credential.user.getIdToken();
                const res: Response = await fetch(
                    "/api/signin", {
                        method: "POST",
                        body: JSON.stringify({ token }),
                    },
                );
                const result = await res.json();
                console.log(result);
            } catch (err) {
                window.console.error(err.message);
            } finally {
                setSubmitting(false);
            }
        })();
    }, []);
    return (
      <div>
        <header className="App-header">
          <button
              className="btn btn-outline-primary"
              disabled={submitting}
              onClick={onClickSignIn}>Sign in</button>
        </header>
      </div>
    );
};

export default Index;
