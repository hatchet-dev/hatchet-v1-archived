import { MaterialIcon } from "@hatchet-dev/hatchet-components";
import React, { Component, useRef } from "react";
import { StatusContainer, StatusText } from "./styles";

type Props = {
  status_text: string;
  material_icon?: string;
  kind?: "icon" | "color";
  color?: string;
};

const Status: React.FC<Props> = ({
  status_text,
  material_icon,
  kind = "icon",
  color = "red",
}) => {
  if (kind == "color") {
    return (
      <StatusContainer color={color}>
        <StatusText color={color}>{status_text}</StatusText>
      </StatusContainer>
    );
  }
  return (
    <StatusContainer>
      <MaterialIcon className="material-icons">{material_icon}</MaterialIcon>
      <StatusText>{status_text}</StatusText>
    </StatusContainer>
  );
};

export default Status;
