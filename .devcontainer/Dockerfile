# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.177.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version: 1, 1.16, 1.15
#FROM mcr.microsoft.com/vscode/devcontainers/go:0-1.19
FROM infovistadev/tems-bc-go:1.19-45783

ARG USERNAME=vscode
# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

RUN apt-get update \
    && apt install -y --no-install-recommends xdg-utils \
    # Install Oh-my-bash:
    && sudo su $USERNAME -c 'rm -rf $HOME/.oh-my-bash && bash -c "$(curl -fsSL https://raw.githubusercontent.com/ohmybash/oh-my-bash/master/tools/install.sh)"' \
    #
    # Copy localhost's private SSH key into the container
    && echo '\n\
    if [ "$SYNC_LOCALHOST_SSH_KEY" = "true" ] && [ -f "/usr/local/share/ssh-localhost/id_rsa" ]; then\n\
    sudo cp -r /usr/local/share/ssh-localhost/id_rsa $HOME/.ssh/id_rsa\n\
    sudo chown -R $(id -u) $HOME/.ssh/id_rsa\n\
    fi' | tee -a /root/.bashrc /root/.zshrc /home/${USERNAME}/.bashrc >> /home/${USERNAME}/.zshrc \
    # Clean up
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/

# Install Mage to build the code
RUN su $USERNAME -c "git clone https://github.com/magefile/mage /tmp/gotools/mage \
    && cd /tmp/gotools/mage \
    && go run bootstrap.go"

# Select oh-my-bash theme:
RUN sed -i -e 's/OSH_THEME=.*/OSH_THEME="powerline-multiline"/g' /home/${USERNAME}/.bashrc

ENV MAGEFILE_VERBOSE=true
ENV MAGEFILE_ENABLE_COLOR=true
