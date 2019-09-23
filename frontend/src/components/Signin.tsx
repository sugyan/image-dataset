import React, { useCallback, useState } from "react";
import firebase from "firebase/app";
import "firebase/auth";

const Signin: React.FC = () => {
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
                if (res.ok) {
                    window.console.log("OK!");
                } else {
                    throw new Error(res.statusText);
                }
            } catch (err) {
                alert(err.message);
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

export default Signin;
