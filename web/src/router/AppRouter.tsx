import React from "react";
import {BrowserRouter, Route, Switch} from "react-router-dom";

import {routeUrls} from "../configs/routeUrls";
import Login from "../pages/Login/Login";
import Pictures from "../pages/Pictures/Pictures";
import Picture from "../pages/Picture/Picture";
import Analytics from "../pages/Analytics/Analytics";

export default function AppRouter({
  transport,
  openCreatePictureModal,
}: {
  transport: string;
  closeCreatePictureModal: () => void;
  openCreatePictureModal: () => void;
}): JSX.Element {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path={routeUrls.login}>
          <Login transport={transport} />
        </Route>
        <Route exact path={routeUrls.pictures}>
          <Pictures transport={transport} openCreatePictureModal={openCreatePictureModal} />
        </Route>
        <Route exact path={routeUrls.picture.route}>
          <Picture transport={transport} />
        </Route>
        <Route exact path={routeUrls.analytics}>
          <Analytics transport={transport} />
        </Route>
      </Switch>
    </BrowserRouter>
  );
}
