import github from "assets/github.png";
import React from "react";
import { StatusText } from "components/status/styles";
import { GithubRefContainer, GithubImg } from "./styles";

type Props = {
  text: string;
  link?: string;
};

const GithubRef: React.FC<Props> = ({ link, text }) => {
  return (
    <GithubRefContainer
      onClick={() => {
        link && window.open(link);
      }}
    >
      <GithubImg src={github} />
      <StatusText>{text}</StatusText>
    </GithubRefContainer>
  );
};

export default GithubRef;
