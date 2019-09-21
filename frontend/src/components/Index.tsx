import React from "react";
import { Link } from "react-router-dom";

const Index: React.FC = () => {
    return (
      <div>
        <header className="App-header">
          <Link
              className="App-link"
              to="/signin"
          >
            Sign in
          </Link>
        </header>
      </div>
    );
};

export default Index;
