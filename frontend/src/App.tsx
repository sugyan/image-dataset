import React from "react";
import { BrowserRouter, Route } from "react-router-dom";

import "./App.css";
import Index from "./components/Index";
import Signin from "./components/Signin";

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <div className="App">
                <Route path="/" exact component={Index}></Route>
                <Route path="/signin" exact component={Signin}></Route>
            </div>
        </BrowserRouter>
    );
};

export default App;
