import { MaterialIcon } from "components/globals";
import React from "react";
import { StyledErrorBar, StyledErrorText } from "./styles";

export type Props = {
  text: string;
  severity?: "warning" | "critical";
};

const ErrorBar: React.FC<Props> = ({ text, severity = "critical" }) => {
  const color = severity == "warning" ? "#f5cb42" : "#ff385d";

  return (
    <StyledErrorBar color={color}>
      <MaterialIcon className="material-icons">error</MaterialIcon>
      <StyledErrorText color={color}>{text}</StyledErrorText>
    </StyledErrorBar>
  );
};

export default ErrorBar;
