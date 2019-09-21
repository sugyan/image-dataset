import React from "react";
import { BrowserRouter, Route } from "react-router-dom";

import "./App.css";
import Index from "./components/Index";

const App: React.FC = () => {
    return (
      <BrowserRouter>
        <div className="App">
          <Route path="/" exact component={Index}></Route>
        </div>
      </BrowserRouter>
    );
};

export default App;
