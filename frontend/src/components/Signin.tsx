import React, { useCallback, useState } from "react";
import {
    Button,
    Theme, makeStyles, createStyles,
} from "@material-ui/core";
import firebase from "firebase/app";
import "firebase/auth";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        button: {
            padding: theme.spacing(2),
        },
    }),
);

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
                    window.location.replace("/");
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
    const classes = useStyles();
    return (
      <div>
        <Button
            className={classes.button}
            color="primary"
            disabled={submitting}
            onClick={onClickSignIn}>
          Sign in
        </Button>
      </div>
    );
};

export default Signin;
