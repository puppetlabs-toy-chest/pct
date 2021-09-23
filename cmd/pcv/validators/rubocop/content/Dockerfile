FROM puppet/puppet-agent
RUN /opt/puppetlabs/puppet/bin/gem install rubocop --no-document
ADD rubocop_ruby* /
RUN ln -s /rubocop_ruby$(/opt/puppetlabs/bin/puppet --version | cut -d '.' -f 1).yml /rubocop.yml
ENTRYPOINT [ "/opt/puppetlabs/puppet/bin/rubocop" ]
CMD [ "--config", "/rubocop.yml", "/module" ]
