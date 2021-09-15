import {Button, Container, Input, Typography} from "@material-ui/core";
import React, {ChangeEvent, useState} from "react";
import {Redirect, useHistory} from "react-router-dom";

import {api} from "../../api/api";
import {routeUrls} from "../../configs/routeUrls";
import {getAccessToken, setAccessToken} from "../../shared/helpers";
import {LoginMode} from "./Login.interface";

export default function Login({transport}: {transport: string}): JSX.Element {
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
      const uri = mode === LoginMode.LOGIN ? "login" : "signup";
      const {data} = await api.post(`/${transport}/${uri}`, {username, password});
      setAccessToken(data["access_token"]);
      history.replace(routeUrls.pictures);
    } catch (err: any) {
      setError(err.message || err.toString());
    }
  };

  return token ? (
    <Redirect to={routeUrls.pictures} />
  ) : (
    <Container style={{marginTop: "30%"}}>
      <Typography variant="h3" component="h3" style={{marginBottom: 50}}>
        Go Realworld
      </Typography>
      <div
        style={{
          marginTop: 24,
          display: "flex",
          margin: "auto",
          flexDirection: "column",
          width: "300px",
        }}
      >
        <Input value={username} onChange={handleUsernameChange} style={{marginBottom: 12}} />
        <Input
          type="password"
          value={password}
          onChange={handlePasswordChange}
          style={{marginBottom: 12}}
        />
        <Button
          color="primary"
          variant="contained"
          onClick={handleSubmitClick}
          style={{marginBottom: 12, fontWeight: 700}}
        >
          {mode === LoginMode.LOGIN ? "Login" : "Sign Up"}
        </Button>
      </div>
      {error ? <p>{error}</p> : null}
      <div>
        {mode === LoginMode.LOGIN ? (
          <div>
            Don't have account?{" "}
            <span onClick={handleModeClick} style={{cursor: "pointer", color: "#00f"}}>
              SignUp
            </span>
          </div>
        ) : (
          <div>
            Alreeady member?{" "}
            <span onClick={handleModeClick} style={{cursor: "pointer", color: "#00f"}}>
              Login
            </span>
          </div>
        )}
      </div>
    </Container>
  );
}
