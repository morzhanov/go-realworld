import React, {useState} from "react";
import {QueryClientProvider, QueryClient} from "react-query";

import "./App.css";

import AppRouter from "../../router/AppRouter";
import TransportSwitch from "../../components/TransportSwitch/TransportSwitch";
import CreatePictureModal from "../../components/CreatePictureModal/CreatePictureModal";

const queryClient = new QueryClient();

export default function App() {
  const [transport, changeTransport] = useState("events");
  const [createPictureModal, setOpenCreatePictureModal] = React.useState(false);
  const handleTransportChange = (e: any): void => {
    changeTransport(e.target.value);
  };
  const openCreatePictureModal = () => setOpenCreatePictureModal(true);
  const closeCreatePictureModal = () => setOpenCreatePictureModal(false);
  return (
    <div className="App">
      <QueryClientProvider client={queryClient}>
        <AppRouter
          transport={transport}
          openCreatePictureModal={openCreatePictureModal}
          closeCreatePictureModal={closeCreatePictureModal}
        />
      </QueryClientProvider>
      <div style={{position: "absolute", bottom: "0", left: "0"}}>
        <TransportSwitch value={transport} handleChange={handleTransportChange} />
      </div>
      <CreatePictureModal
        open={createPictureModal}
        handleClose={closeCreatePictureModal}
        transport={transport}
      />
    </div>
  );
}
