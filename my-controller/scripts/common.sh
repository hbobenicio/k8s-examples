#!/bin/bash
set -eu -o pipefail

ANSI_ESC="\x1b"
ANSI_RESET="${ANSI_ESC}[0m"

function ansi_color_256() {
    local color_code="$1"
    echo "${ANSI_ESC}[38;5;${color_code}m"
}

# https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
LOG_POC_COLOR=$(ansi_color_256 214)
LOG_INFO_COLOR=$(ansi_color_256 121)
LOG_WARN_COLOR=$(ansi_color_256 166)
LOG_ERROR_TAG_COLOR="${ANSI_ESC}[1;31m"
LOG_ERROR_MSG_COLOR=$(ansi_color_256 160)
LOG_MSG_COLOR=$(ansi_color_256 248)  # 187

function log_info() {
    local log_level="info"
    echo -e "${LOG_POC_COLOR}poc: ${LOG_INFO_COLOR}${log_level}:$LOG_MSG_COLOR" $@ "$ANSI_RESET" >&2
}

function log_warn() {
    local log_level="warn"
    echo -e "${LOG_POC_COLOR}poc: ${LOG_WARN_COLOR}${log_level}:$LOG_MSG_COLOR" $@ "$ANSI_RESET" >&2
}

function log_error() {
    local log_level="error"
    echo -e "${LOG_POC_COLOR}poc: ${LOG_ERROR_TAG_COLOR}${log_level}:" $@ "$ANSI_RESET" >&2
}

function confirm_kubectl_context() {
    log_warn "ATENÇÃO! este script utiliza comandos kubectl considerando o contexto ativo atualmente."
    log_warn "Atente-se de configurar o contexto correto para o cluster desejado como alvo da operacao."
    CURRENT_CONTEXT="$(kubectl config current-context)"
    log_warn "contexto kubectl atual: ${LOG_WARN_COLOR}$CURRENT_CONTEXT${ANSI_RESET}"
    read -p "deseja prosseguir com a poc para este cluster? [y/n] " OPTION
    case "$OPTION" in
        n|N)
            log_info "execução cancelada pelo usuário. nada a fazer."
            exit 0
            ;;
        y|Y)
            log_info "prosseguindo com a execucao..."
            ;;
        *)
            log_warn "opcao invalida. opcoes_validas=['y', 'Y', 'n', 'N']."
            log_warn "cancelando a execução por via das dúvidas."
            exit 1
    esac
}
