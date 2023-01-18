import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  P,
  StandardButton,
  AppWrapper,
  ErrorBar,
  SectionAreaWithLogo,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";
import useQueryParam from "shared/hooks/usequeryparam";
import { stringify } from "qs";

const AcceptOrganizationInviteView: React.FunctionComponent = () => {
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");
  const history = useHistory();
  const query = useQueryParam();

  const tokenId = query.get("invite_id");
  const token = query.get("token");
  const orgName = query.get("org_name");
  const inviterAddress = query.get("inviter_address");

  useEffect(() => {
    if (!tokenId || !token || !orgName || !inviterAddress) {
      history.push("/");
    }
  }, [tokenId, token, orgName, inviterAddress]);

  const acceptMutation = useMutation({
    mutationKey: ["accept_invite"],
    mutationFn: () => {
      return api.acceptOrgMemberInvite(tokenId, token);
    },
    onSuccess: (data) => {
      setSuccess(true);
      history.push("/");
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    if (!success) {
      acceptMutation.mutate();
    }
  };

  const renderContents = () => {
    return (
      <SectionAreaWithLogo width="400px">
        <HorizontalSpacer spacepixels={18} />
        <FlexColCenter>
          <H2>Accept Invite?</H2>
        </FlexColCenter>
        <HorizontalSpacer
          spacepixels={10}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={30} />
        {success && (
          <P>
            You have successfully accepted this invite. You will be
            automatically redirected to the Hatchet dashboard.
          </P>
        )}
        {!success && (
          <P>
            You have received an invite from {inviterAddress} to join {orgName}.
            Would you like to accept this invite?
          </P>
        )}
        <HorizontalSpacer spacepixels={30} />
        {err && <ErrorBar text={err} />}
        {err && <HorizontalSpacer spacepixels={12} />}
        {success && (
          <FlexRowRight>
            <StandardButton
              label="Go to Dashboard"
              material_icon="chevron_right"
              icon_side="right"
              on_click={() => {
                history.push("/");
              }}
              margin={"0"}
            />
          </FlexRowRight>
        )}
        {!success && (
          <FlexRowRight>
            <StandardButton
              label="Accept Invite"
              material_icon="check"
              icon_side="right"
              on_click={() => {
                submit();
              }}
              margin={"0"}
            />
          </FlexRowRight>
        )}
      </SectionAreaWithLogo>
    );
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>{renderContents()}</FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default AcceptOrganizationInviteView;
