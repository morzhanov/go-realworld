import React, {ChangeEvent, useState} from "react";
import {makeStyles, Theme, createStyles} from "@material-ui/core/styles";
import Modal from "@material-ui/core/Modal";
import {Button, Input} from "@material-ui/core";

import {api} from "../../api/api";
import {getAuthorization} from "../../shared/helpers";

function rand() {
  return Math.round(Math.random() * 20) - 10;
}

function getModalStyle() {
  const top = 50 + rand();
  const left = 50 + rand();

  return {
    top: `${top}%`,
    left: `${left}%`,
    transform: `translate(-${top}%, -${left}%)`,
  };
}

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    paper: {
      position: "absolute",
      width: 400,
      backgroundColor: theme.palette.background.paper,
      border: "2px solid #000",
      boxShadow: theme.shadows[5],
      padding: theme.spacing(2, 4, 3),
    },
  })
);

export default function CreatePictureModal({
  open,
  handleClose,
  transport,
}: {
  open: boolean;
  handleClose: () => void;
  transport: string;
}) {
  const [title, setTitle] = useState("");
  const [base64, setBase64] = useState("");
  const classes = useStyles();
  // getModalStyle is not a pure function, we roll the style only on the first render
  const [modalStyle] = React.useState(getModalStyle);
  const handleSubmitClick = async (): Promise<void> => {
    try {
      await api.post(
        `${transport}/pictures`,
        {title, base64},
        {
          headers: {Authorization: getAuthorization()},
        }
      );
      handleClose();
    } catch (err: any) {
      console.log(err);
    }
  };
  const handleTitleChange = (e: ChangeEvent<HTMLInputElement>): void => {
    setTitle(e.target.value);
  };
  const handleBase64Change = (e: ChangeEvent<HTMLInputElement>): void => {
    setBase64(e.target.value);
  };

  return (
    <div>
      <Modal
        open={open}
        onClose={handleClose}
        aria-labelledby="simple-modal-title"
        aria-describedby="simple-modal-description"
      >
        <div style={modalStyle} className={classes.paper}>
          <h2 id="simple-modal-title">Create Picture</h2>
          <form
            style={{
              marginTop: 24,
              display: "flex",
              margin: "auto",
              flexDirection: "column",
              width: "300px",
            }}
          >
            <label>Title</label>
            <Input value={title} onChange={handleTitleChange} style={{marginBottom: 12}} />
            <label>Base64 Image</label>
            <Input value={base64} onChange={handleBase64Change} style={{marginBottom: 12}} />
            <Button
              color="primary"
              variant="contained"
              onClick={handleSubmitClick}
              style={{marginBottom: 12, fontWeight: 700}}
            >
              Create
            </Button>
          </form>
        </div>
      </Modal>
    </div>
  );
}
