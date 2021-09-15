import React from "react";
import {FormControl, RadioGroup, FormControlLabel, Radio} from "@material-ui/core";

export default function TransportSwitch({
  value,
  handleChange,
}: {
  value: string;
  handleChange: (event: React.ChangeEvent<HTMLInputElement>, value: string) => void;
}) {
  return (
    <FormControl component="fieldset" style={{padding: 20}}>
      <RadioGroup
        aria-label="transport"
        name="transport"
        value={value}
        onChange={handleChange}
        style={{display: "flex", flexWrap: "nowrap", flexDirection: "row"}}
      >
        <FormControlLabel value="rest" control={<Radio color="primary" />} label="REST" />
        <FormControlLabel value="grpc" control={<Radio color="primary" />} label="GRPC" />
        <FormControlLabel value="events" control={<Radio color="primary" />} label="Events" />
      </RadioGroup>
    </FormControl>
  );
}
