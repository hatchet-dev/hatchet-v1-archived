import ClipboardJS from "clipboard";
import React, { useEffect, useRef } from "react";
import styled from "styled-components";

type Props = {
  text: string;
  onSuccess?: (e: ClipboardJS.Event) => void;
  onError?: (e: ClipboardJS.Event) => void;
  wrapperProps?: any;
  as?: any;
  children?: React.ReactNode;
};

const CopyToClipboard: React.FC<Props> = ({
  text,
  onSuccess,
  onError,
  wrapperProps,
  as,
  children,
}) => {
  const triggerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const trigger = triggerRef.current;

    if (!trigger) {
      console.error("Couldn't mount clipboardjs on wrapper component");
      return;
    }

    const clipboard = new ClipboardJS(trigger, {
      text: () => text,
    });

    clipboard.on("success", (e) => {
      onSuccess && onSuccess(e);
    });

    onError && clipboard.on("error", onError);

    return () => clipboard.destroy();
  }, []);

  let wrappedProps = {
    ...(wrapperProps || {}),
  };

  wrappedProps["id"] = "tooltip-anchor";

  return (
    <span ref={triggerRef}>
      <DynamicSpanComponent
        id="tooltip-anchor"
        as={as || "span"}
        {...wrappedProps}
      >
        {children}
      </DynamicSpanComponent>
    </span>
  );
};

export default CopyToClipboard;

const DynamicSpanComponent = styled.span``;
