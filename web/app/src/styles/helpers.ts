import { css } from "styled-components";
import { lighten, darken } from "polished";

export const resolveBackgroundColor = css`
  ${props =>
  `background-color: ${props.current ? lighten(0.5, "#5f7481") : "#white"}`}
`;

export const resolveColor = css`
  ${props => `color: ${props.current ? darken(0.25, "#f9aa33") : "#344955"}`}
`;
