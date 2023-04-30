import React from "react";
import { FooterWrapper } from "./styles";

type Props = {
  children?: React.ReactNode;
};

const Footer: React.FunctionComponent<Props> = ({ children }) => {
  return <FooterWrapper>{children}</FooterWrapper>;
};

export default Footer;
