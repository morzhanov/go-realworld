import {Container, Paper} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect, useParams, useHistory} from "react-router-dom";
import React from "react";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken, getAuthorization} from "../../shared/helpers";
import {PictureData} from "../Pictures/Pictures.interface";

export default function Picture({transport}: {transport: string}): JSX.Element {
  const {id: pictureId} = useParams<{id: string}>();
  const token = getAccessToken();
  const history = useHistory();

  const {data: picture, error} = useQuery<PictureData, AxiosError, PictureData, any>(
    `picture:${pictureId}`,
    () =>
      api
        .get(`${transport}/pictures/${pictureId}`, {headers: {Authorization: getAuthorization()}})
        .then((res) => res.data)
  );
  const handleBackClick = () => history.push(routeUrls.pictures);

  return token ? (
    <Container style={{padding: 20}}>
      <div
        style={{
          color: "blue",
          fontWeight: 700,
          fontSize: 14,
          marginBottom: 12,
          textAlign: "left",
          width: "100%",
          display: "block",
        }}
      >
        <span onClick={() => handleBackClick()} style={{cursor: "pointer"}}>
          Back
        </span>
      </div>
      <Paper>
        {picture && (
          <img
            alt="pic"
            src={picture?.base64}
            title={picture?.title}
            style={{width: "calc(100vh - 200px)"}}
          />
        )}
      </Paper>
      {!!error ? <p>{error.message}</p> : null}
    </Container>
  ) : (
    <Redirect to={routeUrls.login} />
  );
}
