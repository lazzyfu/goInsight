import { colorTokens, radiusScale, typographyScale } from './tokens'

export const antdTheme = {
  token: {
    colorPrimary: colorTokens.primary,
    colorInfo: colorTokens.info,
    colorSuccess: colorTokens.success,
    colorWarning: colorTokens.warning,
    colorError: colorTokens.error,
    colorBgLayout: colorTokens.pageBg,
    colorBgContainer: colorTokens.containerBg,
    colorBorder: colorTokens.border,
    colorText: colorTokens.textPrimary,
    colorTextSecondary: colorTokens.textSecondary,
    borderRadius: radiusScale.lg,
    fontSize: typographyScale.body,
  },
  components: {
    Layout: {
      headerHeight: 52,
      headerBg: colorTokens.containerBg,
      bodyBg: colorTokens.pageBg,
      siderBg: '#142331',
      triggerBg: '#142331',
      triggerColor: '#d8e3ea',
    },
    Card: {
      borderRadiusLG: radiusScale.card,
      bodyPadding: 16,
      headerBg: 'transparent',
    },
    Button: {
      borderRadius: radiusScale.lg,
      controlHeight: 40,
      controlHeightLG: 44,
      fontWeight: 500,
    },
    Input: {
      borderRadius: radiusScale.lg,
      controlHeight: 40,
    },
    Select: {
      borderRadius: radiusScale.lg,
      controlHeight: 40,
    },
    Table: {
      headerBg: '#f8fbfc',
      headerColor: colorTokens.textPrimary,
    },
    Tag: {
      borderRadiusSM: 999,
      fontSizeSM: typographyScale.caption,
    },
  },
}
