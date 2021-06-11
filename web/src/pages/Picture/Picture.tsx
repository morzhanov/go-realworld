import {Container, Paper} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect, useParams} from "react-router-dom";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken} from "../../shared/helpers";
import {Picture} from "../Pictures/Pictures.interface";

export default function Picture(): JSX.Element {
  const {id: pictureId} = useParams<{id: string}>();
  const token = getAccessToken();

  const {data: picture, error} = useQuery<Picture, AxiosError, Picture, any>(
    `picture:${pictureId}`,
    () => api.get(`/pictures/${pictureId}`).then((res) => res.data)
  );

  return token ? (
    <Container>
      <Paper>{picture && <img src={picture?.base64} title={picture?.title} />}</Paper>
      {!!error ? <p>{error.message}</p> : null}
    </Container>
  ) : (
    <Redirect to={routeUrls.login} />
  );
}
