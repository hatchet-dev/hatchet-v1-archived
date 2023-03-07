import React, { Component, useRef, useState } from "react";
import { Log, LogContainer } from "./styles";

type Props = {
  logs: string[];
};

const Logs: React.FC<Props> = ({ logs }) => {
  const [levels, setLevels] = useState(["info", "error"]);

  const getParsedLog = (log: string) => {
    try {
      const parsedLog = JSON.parse(log);

      if (levels.includes(parsedLog.level)) {
        return `[${parsedLog.level}] ${parsedLog.message}`;
      }

      return null;
    } catch {
      return log;
    }
  };

  return (
    <LogContainer>
      {logs.map((log) => {
        const parsedLog = getParsedLog(log);

        if (parsedLog) {
          return <Log>{parsedLog}</Log>;
        }
      })}
    </LogContainer>
  );
};

export default Logs;
