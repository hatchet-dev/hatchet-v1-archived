import {
  FlexRowLeft,
  MaterialIcon,
  SmallSpan,
  P,
} from "@hatchet-dev/hatchet-components";
import React, { Component, useRef, useState } from "react";
import { LinkButtonContainer } from "./styles";

type Props = {
  text: string;
  link: string;
};

const LinkButton: React.FC<Props> = ({ text, link }) => {
  return (
    <LinkButtonContainer to={link}>
      <SmallSpan>{text}</SmallSpan>
      <MaterialIcon className="material-icons">north_east</MaterialIcon>
    </LinkButtonContainer>
  );
};

export default LinkButton;
