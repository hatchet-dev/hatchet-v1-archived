import React, { Component, useRef } from "react";
import { Log, LogContainer } from "./styles";

type Props = {
  logs: string[];
};

const Logs: React.FC<Props> = ({ logs }) => {
  return (
    <LogContainer>
      {logs.map((log) => {
        return <Log>{log}</Log>;
      })}
    </LogContainer>
  );
};

export default Logs;
