import { Relative, Span } from "components/globals";
import React, { useEffect, useRef, useState } from "react";
import { StyledInput } from "./styles";

export type Selection = {
  label: string;
  value: string;
  icon?: string;
};

export type Props = {
  placeholder: string;
};

const TextInput: React.FC<Props> = ({ placeholder }) => {
  const [selection, setSelection] = useState<Selection>();

  return <StyledInput placeholder={placeholder} />;
};

export default TextInput;
