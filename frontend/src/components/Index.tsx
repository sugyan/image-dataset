import React, { useCallback } from "react";
import firebase from "firebase/app";
import "firebase/auth";

const Index: React.FC = () => {
    const onClickSignIn = useCallback(() => {
        const provider = new firebase.auth.GoogleAuthProvider();
        firebase.auth().signInWithPopup(provider).then((credential: firebase.auth.UserCredential) => {
            if (!credential.user) {
                throw new Error("no credential user");
            }
            credential.user.getIdToken().then((token: string) => {
                window.console.log(token);
            }).catch((err: Error) => {
                window.console.error(err.message);
            });
        }).catch((err: Error) => {
            window.console.error(err.message);
        });
    }, []);
    return (
      <div>
        <header className="App-header">
          <button
              className="btn btn-outline-primary"
              onClick={onClickSignIn}>Sign in</button>
        </header>
      </div>
    );
};

export default Index;
