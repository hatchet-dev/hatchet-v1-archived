import ClipboardJS from "clipboard";
import React, { useEffect, useRef } from "react";
import styled from "styled-components";

type Props = {
  text: string;
  onSuccess?: (e: ClipboardJS.Event) => void;
  onError?: (e: ClipboardJS.Event) => void;
  wrapperProps?: any;
  as?: any;
};

const CopyToClipboard: React.FC<Props> = (props) => {
  const triggerRef = useRef();

  useEffect(() => {
    const trigger = triggerRef.current;

    if (!trigger) {
      console.error("Couldn't mount clipboardjs on wrapper component");
      return;
    }

    const clipboard = new ClipboardJS(trigger, {
      text: () => props.text,
    });

    clipboard.on("success", (e) => {
      props.onSuccess && props.onSuccess(e);
    });

    props.onError && clipboard.on("error", props.onError);

    return () => clipboard.destroy();
  }, []);

  let wrappedProps = {
    ...(props.wrapperProps || {}),
  };

  wrappedProps["id"] = "tooltip-anchor";

  return (
    <span ref={triggerRef}>
      <DynamicSpanComponent
        id="tooltip-anchor"
        as={props.as || "span"}
        {...wrappedProps}
      >
        {props.children}
      </DynamicSpanComponent>
    </span>
  );
};

export default CopyToClipboard;

const DynamicSpanComponent = styled.span``;
