import {Container, GridList, GridListTile, makeStyles} from "@material-ui/core";
import {AxiosError} from "axios";
import {useQuery} from "react-query";
import {Redirect} from "react-router-dom";
import {api} from "../../api/api";

import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken} from "../../shared/helpers";
import {Picture} from "./Pictures.interface";

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

export default function Pictures(): JSX.Element {
  const classes = useStyles();
  const token = getAccessToken();

  const {data: pictures, error} = useQuery<Picture[], AxiosError, Picture[], any>("pictures", () =>
    api.get("/pictures").then((res) => res.data)
  );

  return token ? (
    <Container>
      <GridList cellHeight={160} className={classes.gridList} cols={3}>
        {pictures?.map((pic: Picture) => (
          <GridListTile key={pic.title} cols={1}>
            <img src={pic.base64} alt={pic.title} />
          </GridListTile>
        ))}
      </GridList>
      {!!error ? <p>{error.message}</p> : null}
    </Container>
  ) : (
    <Redirect to={routeUrls.login} />
  );
}
