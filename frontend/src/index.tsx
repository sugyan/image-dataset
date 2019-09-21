import React from "react";
import ReactDOM from "react-dom";
import firebase from "firebase/app";
import "bootstrap/dist/css/bootstrap.css";

import "./index.css";
import App from "./App";

const firebaseConfig = {
    apiKey: "AIzaSyD9GAO1u4nL4mU0F0nouOr6R2MNE9hS-ks",
    authDomain: "idol-face-dataset.firebaseapp.com",
    databaseURL: "https://idol-face-dataset.firebaseio.com",
    projectId: "idol-face-dataset",
    storageBucket: "idol-face-dataset.appspot.com",
    messagingSenderId: "1007816730072",
    appId: "1:1007816730072:web:b880d26b1c15d234c57ead"
};
firebase.initializeApp(firebaseConfig);

ReactDOM.render(<App />, document.getElementById("root"));
