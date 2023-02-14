import { MaterialIcon } from "@hatchet-dev/hatchet-components";
import React, { Component, useRef } from "react";
import { StatusContainer, StatusText } from "./styles";

type Props = {
  status_text: string;
  material_icon: string;
};

const Status: React.FC<Props> = ({ status_text, material_icon }) => {
  return (
    <StatusContainer>
      <MaterialIcon className="material-icons">{material_icon}</MaterialIcon>
      <StatusText>{status_text}</StatusText>
    </StatusContainer>
  );
};

export default Status;
