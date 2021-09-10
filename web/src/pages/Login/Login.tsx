import {Button, Container, Input, Typography} from "@material-ui/core";
import {ChangeEvent, useState} from "react";
import {Redirect, useHistory} from "react-router-dom";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken, setAccessToken} from "../../shared/helpers";
import {LoginMode} from "./Login.interface";

export default function Login(): JSX.Element {
  const token = getAccessToken();
  const history = useHistory();
  const [mode, setMode] = useState<LoginMode>(LoginMode.LOGIN);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleUsernameChange = (e: ChangeEvent<HTMLInputElement>): void => {
    setError("");
    setUsername(e.target.value);
  };
  const handlePasswordChange = (e: ChangeEvent<HTMLInputElement>): void => {
    setError("");
    setPassword(e.target.value);
  };
  const handleModeClick = (): void => {
    setError("");
    setMode(mode === LoginMode.LOGIN ? LoginMode.SIGNUP : LoginMode.LOGIN);
  };
  const handleSubmitClick = async (): Promise<void> => {
    try {
      const uri = mode === LoginMode.LOGIN ? "/login" : "/signup";
      const {
        data: {accessToken},
      } = await api.post(uri, {username, password}).then((res: any) => res.data);
      setAccessToken(accessToken);
      history.replace(routeUrls.pictures);
    } catch (err: any) {
      setError(err);
    }
  };

  return token ? (
    <Redirect to={routeUrls.pictures} />
  ) : (
    <Container>
      <Typography variant="h1" component="h1">
        {mode === LoginMode.LOGIN ? "Login" : "Sign Up"}
      </Typography>
      <div style={{marginTop: 24}}>
        <Input value={username} onChange={handleUsernameChange} />
        <Input value={password} onChange={handlePasswordChange} />
        <Button color="primary" onClick={handleSubmitClick}>
          {mode === LoginMode.LOGIN ? "Login" : "Sign Up"}
        </Button>
      </div>
      {error ? <p>{error}</p> : null}
      <p>
        {mode === LoginMode.LOGIN ? (
          <div>
            Don't have account? <span onClick={handleModeClick}>SignUp</span>
          </div>
        ) : (
          <div>
            Alreeady member? <span onClick={handleModeClick}>Login</span>
          </div>
        )}
      </p>
    </Container>
  );
}
