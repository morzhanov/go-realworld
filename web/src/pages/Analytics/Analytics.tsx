import {Container} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect} from "react-router-dom";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken} from "../../shared/helpers";

export default function Analytics(): JSX.Element {
  const token = getAccessToken();

  const {data: analysis, error} = useQuery<any, AxiosError, any, any>("analysis", () =>
    api.get(`/analysis`).then((res) => res.data)
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
