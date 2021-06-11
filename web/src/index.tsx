import React from "react";
import {createBrowserHistory, History} from "history";
import ReactDOM from "react-dom";

import "./index.css";

import App from "./containers/App/App";

const history: History = createBrowserHistory();
ReactDOM.render(<App history={history} />, document.getElementById("root"));
