import React from "react";
import {History} from "history";
import {Route, Router, Switch} from "react-router";

import {routeUrls} from "../configs/routeUrls";

const Login = React.lazy(() => import("../pages/Login/Login"));
const Pictures = React.lazy(() => import("../pages/Pictures/Pictures"));
const Picture = React.lazy(() => import("../pages/Picture/Picture"));
const Analysis = React.lazy(() => import("../pages/Analytics/Analytics"));

export default function AppRouter({history}: {history: History}): JSX.Element {
  return (
    <Router history={history}>
      <React.Suspense fallback={<div>Loading...</div>}>
        <Switch>
          <Route exact path={routeUrls.login}>
            <Login />
          </Route>
          <Route exact path={routeUrls.pictures}>
            <Pictures />
          </Route>
          <Route exact path={routeUrls.picture}>
            <Picture />
          </Route>
          <Route exact path={routeUrls.analysis}>
            <Analysis />
          </Route>
        </Switch>
      </React.Suspense>
    </Router>
  );
}
