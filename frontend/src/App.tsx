import React, { useEffect } from "react";
import { BrowserRouter, Route } from "react-router-dom";
import { AppBar, Toolbar, Typography, Button, createStyles, makeStyles } from "@material-ui/core";

import Index from "./components/Index";
import Signin from "./components/Signin";
import Signout from "./components/Signout";

interface UserInfo {
    displayName: string
    email: string
    rawId: string
}

const useStyles = makeStyles(() => {
    return createStyles({
        title: {
            flexGrow: 1
        },
    });
});

const App: React.FC = () => {
    useEffect(() => {
        (async () => {
            try {
                const res = await fetch("/api/userinfo");
                const user: UserInfo = await res.json();
                window.console.log(user);
            } catch (err) {
                window.console.error(err);
            }
        })();
    }, []);
    const classes = useStyles();
    return (
      <div>
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" className={classes.title}>
              Dataset
            </Typography>
            <Button color="inherit">Login</Button>
          </Toolbar>
        </AppBar>
        <BrowserRouter>
          <Route path="/" exact component={Index}></Route>
          <Route path="/signin" exact component={Signin}></Route>
          <Route path="/signout" exact component={Signout}></Route>
        </BrowserRouter>
      </div>
    );
};

export default App;
