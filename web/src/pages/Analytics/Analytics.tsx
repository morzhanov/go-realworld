import {Container} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect} from "react-router-dom";
import React from "react";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken, getAuthorization} from "../../shared/helpers";

export default function Analytics({transport}: {transport: string}): JSX.Element {
  const token = getAccessToken();

  const {data: analysis, error} = useQuery<any, AxiosError, any, any>("analysis", () =>
    api
      .get(`${transport}/analytics`, {headers: {Authorization: getAuthorization()}})
      .then((res) => res.data)
  );

  return token ? (
    <Container>
      {analysis && analysis}
      {!!error ? <p>{error.message}</p> : null}
    </Container>
  ) : (
    <Redirect to={routeUrls.login} />
  );
}
