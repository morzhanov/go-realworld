import {Container, GridList, GridListTile, makeStyles} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect, useHistory} from "react-router-dom";
import React from "react";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken, getAuthorization} from "../../shared/helpers";
import {PictureData} from "./Pictures.interface";

const useStyles = makeStyles(() => ({
  root: {
    display: "flex",
    flexWrap: "wrap",
    justifyContent: "space-around",
    overflow: "hidden",
  },
  gridList: {
    width: 500,
    height: 450,
  },
}));

export default function Pictures({
  transport,
  openCreatePictureModal,
}: {
  transport: string;
  openCreatePictureModal: () => void;
}): JSX.Element {
  const classes = useStyles();
  const token = getAccessToken();
  const history = useHistory();

  const {data: pictures, error} = useQuery<PictureData[], AxiosError, PictureData[], any>(
    "pictures",
    () =>
      api
        .get(`${transport}/pictures`, {headers: {Authorization: getAuthorization()}})
        .then((res) => res.data)
  );
  const handlePicClick = (id: string) => history.push(routeUrls.picture.link(id));
  const handleCreatePictureClick = () => {
    openCreatePictureModal();
  };
  const handleAnalyticsClick = () => history.push(routeUrls.analytics);

  return token ? (
    pictures?.length ? (
      <Container style={{padding: 20, height: "100%"}}>
        <div style={{width: "100%", marginBottom: 20, padding: 10, display: "flex"}}>
          <span
            onClick={() => handleCreatePictureClick()}
            style={{
              cursor: "pointer",
              fontWeight: 700,
              fontSize: 14,
              color: "blue",
              marginRight: 32,
            }}
          >
            Create Picture
          </span>
          <span
            onClick={() => handleAnalyticsClick()}
            style={{cursor: "pointer", fontWeight: 700, fontSize: 14, color: "blue"}}
          >
            Analytics
          </span>
        </div>
        <GridList
          cellHeight={200}
          spacing={12}
          className={classes.gridList}
          cols={3}
          style={{padding: 0, margin: "auto", height: "calc(100% - 140px)", width: "80%"}}
        >
          {pictures?.map((pic: PictureData) => (
            <GridListTile key={pic.id} cols={1}>
              <img
                src={pic.base64}
                alt={pic.title}
                style={{cursor: "pointer"}}
                onClick={() => handlePicClick(pic.id)}
              />
            </GridListTile>
          ))}
        </GridList>
        {!!error ? <p>{error.message}</p> : null}
      </Container>
    ) : (
      <div />
    )
  ) : (
    <Redirect to={routeUrls.login} />
  );
}
