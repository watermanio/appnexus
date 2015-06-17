package appnexus

import (
	"fmt"
	"net/http"
)

// MemberService handles all requests to the member service API
type MemberService struct {
	*Response
	client *Client
}

// Member usually means the top level AppNexus account to deal with
type Member struct {
	ID                                 int           `json:"id"`
	Name                               string        `json:"name"`
	WhitelabelSupportEmail             string        `json:"whitelabel_support_email"`
	State                              string        `json:"state"`
	NoResellingPriority                int           `json:"no_reselling_priority"`
	EntityType                         string        `json:"entity_type"`
	ResellingExposure                  string        `json:"reselling_exposure"`
	ResellingExposedOn                 string        `json:"reselling_exposed_on"`
	LastModified                       string        `json:"last_modified"`
	Timezone                           string        `json:"timezone"`
	DefaultCurrency                    string        `json:"default_currency"`
	UseInsertionOrders                 bool          `json:"use_insertion_orders"`
	ExposeOptimizationLevers           bool          `json:"expose_optimization_levers"`
	DefaultOptimizationVersion         int           `json:"default_optimization_version"`
	DailyImpsVerified                  int           `json:"daily_imps_verified"`
	DailyImpsSelfAudited               int           `json:"daily_imps_self_audited"`
	DailyImpsUnaudited                 int           `json:"daily_imps_unaudited"`
	AllowNonCpmPayment                 bool          `json:"allow_non_cpm_payment"`
	DefaultAllowCpc                    bool          `json:"default_allow_cpc"`
	DefaultAllowCpa                    bool          `json:"default_allow_cpa"`
	DefaltCurrency                     string        `json:"default_currency,omitempty"`
	DefaultCampaignTrust               string        `json:"default_campaign_trust"`
	DefaultCampaignAllowUnaudited      bool          `json:"default_campaign_allow_unaudited"`
	ContractAllowsUnaudited            bool          `json:"contract_allows_unaudited"`
	EnableFacebook                     bool          `json:"enable_facebook"`
	ReportingDecimalType               string        `json:"reporting_decimal_type"`
	EnableClickAndImpTrackers          bool          `json:"enable_click_and_imp_trackers"`
	DefaultAdProfileID                 int           `json:"default_ad_profile_id"`
	BuyerCreditLimit                   int           `json:"buyer_credit_limit"`
	PlatformExposure                   string        `json:"platform_exposure"`
	ContactEmail                       string        `json:"contact_email"`
	AllowAdProfileOverride             bool          `json:"allow_ad_profile_override"`
	ShortName                          string        `json:"short_name"`
	ExposeEapEcpPlacementSettings      bool          `json:"expose_eap_ecp_placement_settings"`
	DefaultExternalAudit               bool          `json:"default_external_audit"`
	PluginsEnabled                     bool          `json:"plugins_enabled"`
	DefaultPlacementID                 int           `json:"default_placement_id"`
	SellerRevsharePct                  int           `json:"seller_revshare_pct"`
	Dongle                             string        `json:"dongle"`
	AuditNotifyEmail                   string        `json:"audit_notify_email"`
	VisibilityProfileID                int           `json:"visibility_profile_id"`
	PopsEnabledUI                      bool          `json:"pops_enabled_UI"`
	AllowPriorityAudit                 bool          `json:"allow_priority_audit"`
	DefaultAcceptDataProviderUsersync  bool          `json:"default_accept_data_provider_usersync"`
	DefaultAcceptDemandPartnerUsersync bool          `json:"default_accept_demand_partner_usersync"`
	DefaultAcceptSupplyPartnerUsersync bool          `json:"default_accept_supply_partner_usersync"`
	DomainBlacklistEmail               string        `json:"domain_blacklist_email"`
	RequireFacebookPreaudit            bool          `json:"require_facebook_preaudit"`
	PitbullSegmentID                   int           `json:"pitbull_segment_id"`
	PitbullSegmentValue                int           `json:"pitbull_segment_value"`
	Description                        string        `json:"description"`
	SherlockNotifyEmail                string        `json:"sherlock_notify_email"`
	DefaultContentRetrievalTimeoutMs   int           `json:"default_content_retrieval_timeout_ms"`
	DefaultEnableForMediation          bool          `json:"default_enable_for_mediation"`
	PrioritizeMargin                   bool          `json:"prioritize_margin"`
	DealVisibilityProfileID            int           `json:"deal_visibility_profile_id"`
	DeveloperID                        int           `json:"developer_id"`
	DailyBudget                        int           `json:"daily_budget"`
	AccountOwnerUser                   User          `json:"account_owner_user"`
	DefaultCountry                     string        `json:"default_country"`
	ContentCategories                  []interface{} `json:"content_categories"`
	StandardSizes                      []interface{} `json:"standard_sizes"`
}

type memberResponse struct {
	*http.Response
	Obj struct {
		Member  `json:"member"`
		Error   string `json:"error"`
		Status  string `json:"status"`
		Service string `json:"service"`
		Rate    Rate   `json:"dbg_info"`
	} `json:"response"`
}

// Get a member from the Member Service API
func (s *MemberService) Get(memberID int) (*Member, error) {

	path := "member"
	if memberID > 0 {
		path = fmt.Sprintf("%s/%d", path, memberID)
	}

	req, err := s.client.newRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	m := &memberResponse{}
	_, err = s.client.do(req, m)
	if err != nil {
		return nil, err
	}

	member := &m.Obj.Member
	return member, nil
}

// GetDefault AppNexus member object and set the working member in AppNexus.Client
func (s *MemberService) GetDefault() (*Member, error) {
	member, err := s.Get(0)
	if err != nil {
		return nil, err
	}

	s.client.MemberID = member.ID
	return member, nil
}
