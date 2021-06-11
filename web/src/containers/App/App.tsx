import "./App.css";

import AppRouter from "../../router/AppRouter";
import {AppProps} from "./App.interface";

export default function App({history}: AppProps) {
  return (
    <div className="App">
      <AppRouter history={history} />
    </div>
  );
}
