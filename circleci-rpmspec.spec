Name:           signifai-snap-plugin-publisher-signifai
Version:        VERSION
Release:        1%{?dist}
Summary:        Snap Telemetry Agent

License:        Apache
URL:            https://github.com/signifai/snap-plugin-publisher-signifai
Source0:        snap-plugin-publisher-signifai

Requires:       signifai-go >= 1.8.3-el6.1
Requires:       signifai-snap-agent >= 1.2.0-el6.1

%description


%prep
# No prep; already done

%build
# No build; we already did that.

%install
rm -rf $RPM_BUILD_ROOT

mkdir -p $RPM_BUILD_ROOT/opt/signifai/snap/plugins
cp %{SOURCE0} $RPM_BUILD_ROOT/opt/signifai/snap/plugins/snap-plugin-publisher-signifai
%clean


%files
%defattr(-,root,root,-)
/opt/signifai/snap/plugins/snap-plugin-publisher-signifai

%changelog
