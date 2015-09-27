package gsg

type Response struct {
	APIVersion string            `json:"apiVersion"`
	Context    string            `json:"context"`
	ID         string            `json:"id"`
	Method     string            `json:"method"`
	Params     map[string]string `json:"params"`
	Data       `json:"datai, omitempty"`
	Error      `json:"error, omitempty"`
}

func New() *Response {
	r := Response{Params: make(map[string]string), Data: NewData()}
	return &r
}

func NewFromResponse(r http.Response) *Response, err {
	res := New()
	defer r.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	err := json.UnMarshall(body, res)
	return res, err
}

func (r *Response) Copy() *Response {
	nr := Response{
		APIVersion: r.APIVersion,
		Method:     r.Method,
		Params:     r.Params,
	}
	return &nr
}

func (r *Response) Write() ([]byte, error) {
	if r.Error == nil && r.Data == nil {
		err := errors.New("Data and error both set to nil")
		return make([]byte, err)
	}
	js, err := json.Marshal(r)
}

func (r *Response) WriteToResponse(w http.ResponseWriter) error {
	b, err := r.Write()
	w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type Data struct {
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Kind             string            `json:"kind"`
	Fields           string            `json:"fields"`
	Etag             string            `json:"etag"`
	ID               string            `json:"id"`
	Lang             string            `json:"lang"`
	Updated          string            `json:"updated"`
	Deleted          bool              `json:"deleted"`
	CurrentItemCount int               `json:"currentItemCount"`
	ItemsPerPage     int               `json:"itemsPerPage"`
	StartIndex       int               `json:"startIndex"`
	TotalItems       int               `json:"totalItems"`
	PageIndex        int               `json:"pageIndex"`
	TotalPages       int               `json:"totalPages"`
	NextLink         string            `json:"nextLink"`
	PreviousLink     string            `json:"previousLink"`
	SelfLink         string            `json:"selfLink"`
	EditLink         string            `json:"editLink"`
	Items            []json.RawMessage `json:"items"`
	item             int
}

func NewData() *Data {
	d := Data{Items: make([]interface{})}
	return &d
}

func (d *Data) AddField(key string) {
	fs := d.GetFields()
	fs = append(fs, key)
	d.Fields = strings.Join(fs, ",")
}

func (d *Data) GetFields() []string {
	return strings.Split(d.Fields, ",")
}

func (d *Data) AddItem(i interface{}) error {
	js, err := json.Marshall(i)
	d.Items = append(d.Items, js)
	return err
}

func (d *Data) ItemsCount() int {
	return len(d.Items)
}

func (d *Data) NextItem(i interface{}) error {
	count = d.ItemCount()
	if count == d.item {
		return error.New("End of items")
	}
	err := json.UnMarshall(d.Items[d.item], i)
	d.item = d.item + 1
	return err
}

func (d *Data) ResetItems() {
	d.item = 0
}

type Error struct {
	Code    int         `json:"code"`
	Errors  []ErrorItem `json:"errors"`
	Message string      `json:"message"`
}

type ErrorItem struct {
	ExtendedHelper string `json:"extendedHelper"`
	Location       string `json:"location"`
	LocationType   string `json:"locationType"`
	Message        string `json:"message"`
	SendReport     string `json:"sendReport"`
}
