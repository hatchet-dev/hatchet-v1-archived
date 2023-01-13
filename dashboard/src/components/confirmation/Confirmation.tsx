import StandardButton from "components/buttons";
import CopyToClipboard from "components/copytoclipboard";
import {
  FlexCol,
  FlexRow,
  FlexRowRight,
  HorizontalSpacer,
  P,
} from "components/globals";
import TextInput from "components/textinput";
import React, { useState } from "react";

export type Props = {
  confirm_text: string;
  prompt: string;
  button_label: string;
  confirm_text_example?: string;
  confirmed: () => void;
};

const Confirmation: React.FC<Props> = ({
  confirm_text,
  prompt,
  confirm_text_example,
  button_label,
  confirmed,
}) => {
  const [confirmText, setConfirmText] = useState("");

  return (
    <FlexCol>
      <P>{prompt}</P>
      <HorizontalSpacer spacepixels={20} />
      <TextInput
        placeholder={`ex. ${confirm_text_example || "delete"}`}
        type="text"
        width="100%"
        on_change={(val) => {
          setConfirmText(val);
        }}
      />
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label={button_label}
          material_icon="delete"
          icon_side="right"
          on_click={() => {
            confirmed();
          }}
          margin={"0"}
          disabled={confirm_text != confirmText}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default Confirmation;
